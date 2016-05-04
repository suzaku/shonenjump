package main

import (
    "path/filepath"
    "os/user"
)

type Config struct {
    DataDir string
}

func (c Config) getDataPath() string {
    return filepath.Join(c.DataDir, "shonenjump.txt")
}

func getConfig() Config {
    dir := getDefaultDataDir()
    return Config{dir}
}

func getDefaultDataDir() string {
    usr, _ := user.Current()
    dir := filepath.Join(usr.HomeDir, ".local/share/shonenjump")
    return dir
}
