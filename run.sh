#!/bin/zsh

export APP_ENV=local; ent generate ./ent/schema && go install && gosane
