package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var ids_global map[string]int
var topics map[string]Topic

var path_ids_global map[string]int
var steps map[string]Step

type Topic struct {
	Id string
	Name string
	Level string
	Tags []string
	Description string
	Dependencies []string
}

func (t *Topic) Validate() {
	if t.Name == "" || t.Id == "" || t.Level == "" {
		log.Fatalf("Topic %v should have Id, Name, Level fields", t.Id)
	}
	for _, dep := range(t.Dependencies) {
		if ids_global[dep] != 1 {
			log.Fatalf("Failed to resolve dependency id: %v at topic id: %v", dep, t.Id)
		}
	}
}

type Subdomain struct {
	Id string
	Name string
	Topics []Topic
}

func (t *Subdomain) Validate() {
	if t.Name == "" {
		log.Fatalf("Subdomain %v should have Name field", t.Id)
	}
}

type Domain struct {
	Id string
	Name string
	Subdomains []Subdomain
}

func (t *Domain) Validate() {
	if t.Name == "" {
		log.Fatalf("Domain %v should have Name field", t.Id)
	}
}

type Index struct {
	Version string
	Domains []Domain
}


func (c *Index) getTopics(path string) {

    yamlFile, err := ioutil.ReadFile(path)
    if err != nil {
        log.Fatalf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }
}

func (c *Index) ParseTopics() {
	for _, dom := range(c.Domains) {
		if ids_global[dom.Id] == 1 {
			log.Fatalf("Unique id contraint failed on %v domain", dom.Name)
		} else {
			ids_global[dom.Id] = 1
		}
		for _, subdom := range(dom.Subdomains) {
			if ids_global[subdom.Id] == 1 {
				log.Fatalf("Unique Id constraint failed on %v subdomain", subdom.Name)
			} else {
				ids_global[subdom.Id] = 1
			}
			for _, topic := range(subdom.Topics) {
				if ids_global[topic.Id] == 1 {
					log.Fatalf("Unique Id constraint failed on %v topic", topic.Name)
				} else {
					ids_global[topic.Id] = 1
					topics[topic.Id] = topic
				}
			}
		}
	}
}

func (c *Index) Validate(){
	for _, dom := range(c.Domains) {
		dom.Validate()
		for _, subdom := range(dom.Subdomains) {
			subdom.Validate()
			for _, topic := range(subdom.Topics) {
				topic.Validate()
			}
		}
	}

}

type Step struct {
	Id string
	Topic_Id string
}

func (step *Step) Validate(){
	if step.Id == "" || step.Topic_Id == "" {
		log.Fatalf("Step can't have empty Id or Topic_Id fields")
	}
	if ids_global[step.Topic_Id] != 1 {
		log.Fatalf("In step %v non existent topic %v was declared", step.Id, step.Topic_Id)
	}
}

type Chapter struct {
	Id string
	Name string
	Steps []Step
}

type Path struct {
	Version string
	Name string
	Id string
	Level string
	Chapters []Chapter
}

func (path *Path) ReadPath(path_file string){
	yamlFile, err := ioutil.ReadFile(path_file)
    if err != nil {
        log.Fatalf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, path)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }
}

func (path *Path) ReadSteps() {
	for _, chp := range(path.Chapters) {
		if path_ids_global[chp.Id] == 1 {
			log.Fatalf("Unique id contraint failed on %v chapter", chp.Name)
		} else {
			path_ids_global[chp.Id] = 1
		}
		for _, step := range(chp.Steps) {
			if path_ids_global[step.Id] == 1{
				log.Fatalf("Unique id contraint failed on %v step", step.Id)
			} else {
				path_ids_global[step.Id] = 1
				steps[step.Id] = step
			}
		}
	}
}


func (path *Path) Validate() {
	topics_ids_path := make(map[string]bool)
	for _, chp := range(path.Chapters) {
		for _, step := range(chp.Steps) {
			step.Validate()
			for _, dep := range(topics[step.Topic_Id].Dependencies){
				if topics_ids_path[dep] == false {
					log.Fatalf("Dependency %v not met for step.id: %v", dep, step.Id)
				}
			}
			topics_ids_path[step.Topic_Id] = true
		}
	}
}

func load_domain(path_file string){
	index := Index{}
	index.getTopics(path_file)
	index.ParseTopics()
	index.Validate()
}
func main(){
    argsWithoutProg := os.Args[1:]
	cmd := argsWithoutProg[0]
	path_file := argsWithoutProg[1]

	ids_global = make(map[string]int)
	topics = make(map[string]Topic)

	path_ids_global = make(map[string]int)
	steps = make(map[string]Step)

	switch cmd {
		case "domain": {
			load_domain(path_file)
		}
		case "path": {
			load_domain(argsWithoutProg[2])
			path := Path{}
			path.ReadPath(path_file)
			path.ReadSteps()
			path.Validate()
		}
	}
	
}