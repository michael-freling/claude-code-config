.PHONY: install uninstall

CLAUDE_DIR := $(HOME)/.claude
PROJECT_DIR := $(shell pwd)

AGENT_FILES := $(wildcard agents/*.md)
RULE_FILES := $(wildcard rules/*.md)
COMMAND_FILES := $(wildcard commands/*.md)
SKILL_DIRS := $(wildcard skills/*/)

install:
	@echo "Creating symlinks in $(CLAUDE_DIR)..."
	@mkdir -p $(CLAUDE_DIR)/agents $(CLAUDE_DIR)/rules $(CLAUDE_DIR)/commands $(CLAUDE_DIR)/skills
	@for f in $(AGENT_FILES); do \
		ln -sfn $(PROJECT_DIR)/$$f $(CLAUDE_DIR)/$$f; \
	done
	@for f in $(RULE_FILES); do \
		ln -sfn $(PROJECT_DIR)/$$f $(CLAUDE_DIR)/$$f; \
	done
	@for f in $(COMMAND_FILES); do \
		ln -sfn $(PROJECT_DIR)/$$f $(CLAUDE_DIR)/$$f; \
	done
	@for d in $(SKILL_DIRS); do \
		target=$(CLAUDE_DIR)/$$d; \
		target=$${target%/}; \
		if [ -L "$$target" ]; then \
			rm -f "$$target"; \
		elif [ -d "$$target" ]; then \
			rm -rf "$$target"; \
		fi; \
		ln -sfn $(PROJECT_DIR)/$${d%/} "$$target"; \
	done
	@ln -sfn $(PROJECT_DIR)/settings.json $(CLAUDE_DIR)/settings.json
	@echo "Done."

uninstall:
	@echo "Removing symlinks from $(CLAUDE_DIR)..."
	@for f in $(AGENT_FILES); do \
		rm -f $(CLAUDE_DIR)/$$f; \
	done
	@for f in $(RULE_FILES); do \
		rm -f $(CLAUDE_DIR)/$$f; \
	done
	@for f in $(COMMAND_FILES); do \
		rm -f $(CLAUDE_DIR)/$$f; \
	done
	@for d in $(SKILL_DIRS); do \
		target=$(CLAUDE_DIR)/$$d; \
		target=$${target%/}; \
		if [ -L "$$target" ]; then \
			rm -f "$$target"; \
		fi; \
	done
	@rm -f $(CLAUDE_DIR)/settings.json
	@echo "Done."
