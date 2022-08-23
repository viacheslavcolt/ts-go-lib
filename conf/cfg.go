package conf

import (
	"flag"
	"fmt"
)

/*
	Cfg api

	API

	type CfgMap struct { ... }

	NewCfgMap() *CfgMap

	LoadFromFile(cfg *CfgMap, fileName) err
	LoadFromFlag(cfg *CfgMap) err

	GetString(key) string
	GetBool(key) bool
	GetInt(key) int
	Exists(key) bool
*/

const (
	CfgFilePathFlag string = "cfg_path"
)

type Cfg struct {
	entriesCount int
	_map         *_CfgMap
}

func NewCfg() *Cfg {
	return &Cfg{
		entriesCount: 0,
		_map:         _NewCfgMap(),
	}
}

func LoadFromFile(cfg *Cfg, filePath string) error {
	return load(cfg, filePath)
}

func LoadFromFlag(cfg *Cfg) error {
	var (
		filePath string
	)

	flag.StringVar(&filePath, CfgFilePathFlag, "", "configuration file path")
	flag.Parse()

	if len(filePath) == 0 {
		return fmt.Errorf("empty cfg_path option")
	}

	return load(cfg, filePath)
}

func (c *Cfg) GetString(key string) string {
	var (
		obj *_TscObj
	)

	if obj = c._map._Get(key); obj != nil {
		return string(obj.StrVal)
	}

	return ""
}

func (c *Cfg) GetBool(key string) bool {
	var (
		obj *_TscObj
	)

	if obj = c._map._Get(key); obj != nil {
		return obj.Bool
	}

	return false
}

func (c *Cfg) GetInt(key string) int {
	var (
		obj *_TscObj
	)

	if obj = c._map._Get(key); obj != nil {
		return obj.Int
	}

	return 0
}

func (c *Cfg) Exists(key string) bool {
	return c._map._Get(key) != nil
}

// internal

func load(cfg *Cfg, filePath string) error {
	var (
		buf  _Buffer
		root *_TscObj
		err  error
	)

	if err = _InitBuffer(&buf, filePath); err != nil {
		return err
	}

	if root, err = _TscParse(&buf); err != nil {
		return err
	}

	cfg._map._SetRoot(root)

	return nil
}
