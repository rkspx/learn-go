package store

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
)

type NoSuchKeyError struct {
	key string
}

func (err NoSuchKeyError) Error() string {
	return fmt.Sprintf("store: no such key \"%s\"", err.key)
}

type Store struct {
	Data map[string]json.RawMessage
	sync.RWMutex
}

func Open(filename string) (*Store, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var w io.Reader
	if strings.HasSuffix(filename, ".gz") {
		w, err = gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
	} else {
		w = f
	}

	toOpen := make(map[string]string)
	err = json.NewDecoder(w).Decode(&toOpen)
	if err != nil {
		return nil, err
	}

	ks := new(Store)
	ks.Data = make(map[string]json.RawMessage)
	for key, value := range toOpen {
		ks.Data[key] = json.RawMessage(value)
	}
	return ks, nil
}

func Save(ks *Store, filename string) error {
	ks.RLock()
	defer ks.RUnlock()

	toSave := make(map[string]string)
	for key, value := range ks.Data {
		toSave[key] = string(value)
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	if strings.HasSuffix(filename, ".gz") {
		w := gzip.NewWriter(f)
		defer w.Close()
		enc := json.NewEncoder(w)
		enc.SetIndent("", " ")
		return enc.Encode(toSave)
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")
	return enc.Encode(toSave)
}

func (s *Store) Set(key string, value interface{}) error {
	s.Lock()
	defer s.Unlock()
	if s.Data == nil {
		s.Data = make(map[string]json.RawMessage)
	}

	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	s.Data[key] = json.RawMessage(b)
	return nil
}

func (s *Store) Get(key string, value interface{}) error {
	s.RLock()
	defer s.RUnlock()

	b, ok := s.Data[key]
	if !ok {
		return NoSuchKeyError{key}
	}

	return json.Unmarshal(b, &value)
}

func (s *Store) GetAll(re *regexp.Regexp, limit ...int) map[string]json.RawMessage {
	s.RLock()
	defer s.RUnlock()

	results := make(map[string]json.RawMessage)
	for k, v := range s.Data {
		if re == nil || re.MatchString(k) {
			results[k] = v
			if len(limit) > 0 {
				if len(results) >= limit[0] {
					return results
				}
			}
		}
	}

	return results
}

func (s *Store) Keys() []string {
	s.RLock()
	defer s.RUnlock()

	keys := make([]string, len(s.Data))
	i := 0
	for k := range s.Data {
		keys[i] = k
		i++
	}

	return keys
}

func (s *Store) Delete(key string) {
	s.Lock()
	defer s.Unlock()

	delete(s.Data, key)
}
