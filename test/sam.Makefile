# Output helpers
# --------------

TASK_BUILD = echo "🛠️  $@ done"

build-HostingLambda376FBDB5:
	cp ./bootstrap ${ARTIFACTS_DIR}/
	@$(TASK_BUILD)

build-HostingAuthLambdaD197889E:
	cp ./bootstrap ${ARTIFACTS_DIR}/
	@$(TASK_BUILD)
