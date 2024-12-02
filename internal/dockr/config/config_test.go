package config

import (
	"fmt"
	"testing"
)

func TestUltimateConfig(t *testing.T) {
	var configs = []ContainerConfig{
		DefaultConfigDB(),
		DefaultConfigLB(),
		DefaultConfigCache(),
	}

	conf, err := NewContainersConfig(configs...)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(conf.Containers[DB][0])
}

func TestLoadContainersConfig(t *testing.T) {
	t.Run(".json", func(t *testing.T) {
		conf, err := LoadContainersConfig("conf.json")
		if err != nil {
			t.Error(err)
		}

		if len(conf.Containers) != 2 {
			t.Errorf("Expected 2 containers, got %d", len(conf.Containers))
		}

		fmt.Println(conf.Containers[DB][0])
	})

	t.Run(".yaml", func(t *testing.T) {
		conf, err := LoadContainersConfig("conf.yaml")
		if err != nil {
			t.Error(err)
		}

		if len(conf.Containers) != 2 {
			t.Errorf("Expected 2 containers, got %d", len(conf.Containers))
		}

		fmt.Println(conf.Containers[DB][0])
	})

	t.Run(".yml", func(t *testing.T) {
		conf, err := LoadContainersConfig("conf.yml")
		if err != nil {
			t.Error(err)
		}

		if len(conf.Containers) != 2 {
			t.Errorf("Expected 2 containers, got %d", len(conf.Containers))
		}

		fmt.Println(conf.Containers[DB][0])
	})
}

func TestRandString(t *testing.T) {
	four := randString(4)
	ten := randString(10)

	if len(four) != 4 {
		t.Errorf("Expected 4 characters, got %d", len(four))
	}

	if len(ten) != 10 {
		t.Errorf("Expected 10 characters, got %d", len(ten))
	}
}
