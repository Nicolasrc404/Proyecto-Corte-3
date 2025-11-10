package server

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

// redisQueueKey defines the default list used to enqueue tasks.
const redisQueueKey = "alchemy:tasks"

// RedisClient is a minimal RESP client tailored for pushing and popping
// messages from Redis without relying on external dependencies. It only
// implements the handful of commands that the async pipeline requires.
type RedisClient struct {
	addr        string
	dialTimeout time.Duration
}

// NewRedisClient builds a new minimal client for the given address.
func NewRedisClient(addr string) *RedisClient {
	return &RedisClient{addr: addr, dialTimeout: 5 * time.Second}
}

func (c *RedisClient) dial() (net.Conn, *bufio.Reader, error) {
	conn, err := net.DialTimeout("tcp", c.addr, c.dialTimeout)
	if err != nil {
		return nil, nil, err
	}
	return conn, bufio.NewReader(conn), nil
}

// Ping validates connectivity with Redis.
func (c *RedisClient) Ping(ctx context.Context) error {
	conn, reader, err := c.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := writeCommand(conn, "PING"); err != nil {
		return err
	}
	resp, err := parseRESP(ctx, reader)
	if err != nil {
		return err
	}
	if s, ok := resp.(string); ok && strings.EqualFold(s, "PONG") {
		return nil
	}
	return fmt.Errorf("unexpected PING response: %v", resp)
}

// LPUSH prepends a value to the configured list.
func (c *RedisClient) LPUSH(ctx context.Context, key string, value []byte) error {
	conn, reader, err := c.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := writeCommand(conn, "LPUSH", key, string(value)); err != nil {
		return err
	}
	_, err = parseRESP(ctx, reader)
	return err
}

// BRPOP blocks until a value is available for the given key. It honours the
// provided context for cancellation.
func (c *RedisClient) BRPOP(ctx context.Context, key string) ([]byte, error) {
	conn, reader, err := c.dial()
	if err != nil {
		return nil, err
	}

	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			conn.Close()
		case <-done:
		}
	}()

	if err := writeCommand(conn, "BRPOP", key, "0"); err != nil {
		close(done)
		conn.Close()
		return nil, err
	}

	resp, err := parseRESP(ctx, reader)
	close(done)
	conn.Close()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, err
		}
		return nil, err
	}

	arr, ok := resp.([]interface{})
	if !ok || len(arr) != 2 {
		return nil, fmt.Errorf("unexpected BRPOP response: %v", resp)
	}
	payload, ok := arr[1].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected BRPOP payload: %v", arr[1])
	}
	return payload, nil
}

func writeCommand(conn net.Conn, args ...string) error {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "*%d\r\n", len(args))
	for _, arg := range args {
		fmt.Fprintf(&buf, "$%d\r\n%s\r\n", len(arg), arg)
	}
	_, err := conn.Write(buf.Bytes())
	return err
}

func parseRESP(ctx context.Context, reader *bufio.Reader) (interface{}, error) {
	for {
		if ctx != nil {
			select {
			case <-ctx.Done():
				return nil, context.Canceled
			default:
			}
		}
		prefix, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				continue
			}
			return nil, err
		}
		switch prefix {
		case '+':
			line, err := readLine(reader)
			if err != nil {
				return nil, err
			}
			return line, nil
		case '-':
			line, err := readLine(reader)
			if err != nil {
				return nil, err
			}
			return nil, errors.New(line)
		case ':':
			line, err := readLine(reader)
			if err != nil {
				return nil, err
			}
			n, err := strconv.ParseInt(line, 10, 64)
			if err != nil {
				return nil, err
			}
			return n, nil
		case '$':
			line, err := readLine(reader)
			if err != nil {
				return nil, err
			}
			size, err := strconv.Atoi(line)
			if err != nil {
				return nil, err
			}
			if size == -1 {
				return nil, nil
			}
			data := make([]byte, size+2)
			if _, err := io.ReadFull(reader, data); err != nil {
				return nil, err
			}
			return data[:size], nil
		case '*':
			line, err := readLine(reader)
			if err != nil {
				return nil, err
			}
			count, err := strconv.Atoi(line)
			if err != nil {
				return nil, err
			}
			if count == -1 {
				return nil, nil
			}
			arr := make([]interface{}, 0, count)
			for i := 0; i < count; i++ {
				item, err := parseRESP(ctx, reader)
				if err != nil {
					return nil, err
				}
				arr = append(arr, item)
			}
			return arr, nil
		default:
			return nil, fmt.Errorf("unsupported RESP prefix: %q", prefix)
		}
	}
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r"), nil
}
