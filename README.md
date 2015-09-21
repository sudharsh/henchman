THIS REPO HAS BEEN MOVED TO [github.com/apigee/henchman](https://github.com/apigee/henchman)

henchman [![Build Status](https://travis-ci.org/sudharsh/henchman.svg?branch=master)](https://travis-ci.org/sudharsh/henchman)
========

What
----
Henchman is an orchestration and automation tool, inspired by Ansible. Henchman executes a *plan* (a collection of *tasks*) on a given set of machines. 

Why
---
For fun :). Although, python and ruby are awesome as systems languages, I feel Golang fits this 'get shit done' niche more.
Deployments of tools written in Golang are less painful which is a big win in large scale environments where we've had our share of incompatibilities with ruby and python dependencies. 


Current State
-------------
Currently, `henchman` isn't even close to ansible in terms of functionality. For now, it dispatches shell commands on the given hosts and has rudimentary support for variables.
The immediate goals are to nail down the *plan* format and worry about the modules implementation next. Eventually, henchman modules can be written in any language as long as they output to `stdout`, just like ansible.


Building
--------
* Clone this repository
* `make all`
* `bin/henchman -h`
* Refer `samples/plan.yaml`


Contributing
------------
Github PRs away!

