release:
	git tag v$(V)
	@read -p  "Please enter to confirm and push to origin..." && git push origin v$(V)
