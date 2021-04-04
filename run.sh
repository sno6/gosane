#!/bin/zsh

export APP_ENV=local; ent generate ./ent/schema --template ./ent/template && go install && gosane
