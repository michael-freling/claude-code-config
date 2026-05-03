.PHONY: install update uninstall gemini-install

PLUGIN_NAME := michael-config
MARKETPLACE_NAME := michael-freling
PROJECT_DIR := $(shell pwd)
GEMINI_DIR := $(HOME)/.gemini
SKILL_DIRS := $(wildcard skills/*/)

install: claude-install gemini-install

update: claude-update gemini-install

claude-update:
	@CLAUDECODE= claude plugin uninstall $(PLUGIN_NAME)@$(MARKETPLACE_NAME) || true
	@CLAUDECODE= claude plugin install $(PLUGIN_NAME)@$(MARKETPLACE_NAME) --scope user

claude-install:
	@CLAUDECODE= claude plugin marketplace add $(PROJECT_DIR)
	@CLAUDECODE= claude plugin install $(PLUGIN_NAME)@$(MARKETPLACE_NAME) --scope user

gemini-install:
	@echo "Creating symlinks in $(GEMINI_DIR)..."
	@mkdir -p $(GEMINI_DIR)/skills
	@for d in $(SKILL_DIRS); do \
		gtarget=$(GEMINI_DIR)/$$d; \
		gtarget=$${gtarget%/}; \
		if [ -L "$$gtarget" ]; then \
			rm -f "$$gtarget"; \
		elif [ -d "$$gtarget" ]; then \
			rm -rf "$$gtarget"; \
		fi; \
		ln -sfn $(PROJECT_DIR)/$${d%/} "$$gtarget"; \
	done
	@echo "Gemini installation done."

uninstall:
	@CLAUDECODE= claude plugin uninstall $(PLUGIN_NAME)@$(MARKETPLACE_NAME) || true
	@CLAUDECODE= claude plugin marketplace remove $(MARKETPLACE_NAME) || true
	@echo "Removing Gemini symlinks..."
	@for d in $(SKILL_DIRS); do \
		gtarget=$(GEMINI_DIR)/$$d; \
		gtarget=$${gtarget%/}; \
		rm -f "$$gtarget"; \
	done
	@find $(GEMINI_DIR)/skills -xtype l -delete 2>/dev/null || true
	@echo "Done."
