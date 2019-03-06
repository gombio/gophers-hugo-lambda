#!/bin/bash

siege -c1 -t 1s 'https://api.go.gomb.io/captchas POST'

siege -c1 -t 1s 'https://api.go.gomb.io/forms POST {"data":{"title":"Example title","content":"Content"},"captcha":{"id":"abfed3c0-c89d-4e43-a717-65127d8c14e4","secret":"123asd"}}'
