package initialization

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/validator"
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
func GetAllUniqueTags() *[]string {
	tags := make([]string, 0)
	for _, role := range *RoleList {
		tags = append(tags, role.Tags...)
	}
	result := slice.Union(tags)
	return &result
}

func GetRoleByTitle(title string) *Role {
	for _, role := range *RoleList {
		if role.Title == title {
			return &role
		}
	}
	return nil
}

func GetTitleListByTag(tags string) *[]string {
	roles := make([]string, 0)
	//pp.Println(RoleList)
	for _, role := range *RoleList {
		for _, roleTag := range role.Tags {
			if roleTag == tags && !validator.IsEmptyString(role.
				Title) {
				roles = append(roles, role.Title)
			}
		}
	}
	return &roles
}

func GetFirstRoleContentByTitle(title string) (string, error) {
	for _, role := range *RoleList {
		if role.Title == title {
			return role.Content, nil
		}
	}
	return "", errors.New("role not found")
}
