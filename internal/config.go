package internal

import (
	"encoding/json"
	"os"
)

type Config struct {
	Db_Url            string `json:"db_url"`
	Current_User_Name string `json:"current_user_name"`
}

func Read() (Config, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	raw, err := os.ReadFile(dir + "/.gatorconfig.json")
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil

}

func (c *Config) SetUser(name string) error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	c.Current_User_Name = name
	byte_data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dir+"/.gatorconfig.json", byte_data, 0644); err != nil {
		return err
	}

	return nil
}
