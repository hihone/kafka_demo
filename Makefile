.PHONY: test prod
test:
	@sh ./deploy.sh test

prod:
	@sh ./deploy.sh release