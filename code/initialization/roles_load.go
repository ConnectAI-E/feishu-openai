package initialization

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Role struct {
	Title   string   `yaml:"title"`
	Content string   `yaml:"content"`
	Tags    []string `yaml:"tags"`
}

var RoleList *[]Role

// InitRoleList 加载Prompt
func InitRoleList() *[]Role {
	data, err := ioutil.ReadFile("role_list.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &RoleList)
	if err != nil {
		log.Fatal(err)
	}
	return RoleList
}

func GetRoleList() *[]Role {
	return RoleList
}
func GetAllTags() *[]string {
	tags := make([]string, 0)
	for _, role := range *RoleList {
		tags = append(tags, role.Tags...)
	}
	return &tags
}

func GetRoleByTitle(title string) *Role {
	for _, role := range *RoleList {
		if role.Title == title {
			return &role
		}
	}
	return nil
}

func GetRolesByTag(tag string) *[]string {
	roles := make([]string, 0)
	for _, role := range *RoleList {
		for _, roleTag := range role.Tags {
			if roleTag == tag {
				roles = append(roles, role.Title)
			}
		}
	}
	return &roles
}
