package internal

import (
	"encoding/json"
	"os"
)

const DataFileName = "memory.json"

func LoadData() (*AppData, error) {
    f, err := os.Open(DataFileName)
    if err != nil {
        // If the file does not exist, return a default AppData
        if os.IsNotExist(err) {
            return &AppData{NextStoryID: 1}, nil
        }
        return nil, err
    }
    defer f.Close()

    var data AppData
    err = json.NewDecoder(f).Decode(&data)
    if err != nil {
        return nil, err
    }

    return &data, nil
}

func SaveData(data *AppData) error {
    f, err := os.Create(DataFileName)
    if err != nil {
        return err
    }
    defer f.Close()

    enc := json.NewEncoder(f)
    enc.SetIndent("", "  ")
    return enc.Encode(data)
}