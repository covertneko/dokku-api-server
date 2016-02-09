PACKAGE_NAME='github.com/nikelmwann/dokku-api'

all:
	vagrant ssh -c "go install ${PACKAGE_NAME}"
