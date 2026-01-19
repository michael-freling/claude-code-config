.PHONY: install claude-install gemini-install uninstall

CLAUDE_DIR := $(HOME)/.claude
GEMINI_DIR := $(HOME)/.gemini
PROJECT_DIR := $(shell pwd)

AGENT_FILES := $(wildcard agents/*.md)
RULE_FILES := $(wildcard rules/*.md)
COMMAND_FILES := $(wildcard commands/*.md)
SKILL_DIRS := $(wildcard skills/*/)

install: claude-install gemini-install

claude-install:
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
	@echo "Claude installation done."

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
	@echo "Removing symlinks from $(CLAUDE_DIR) and $(GEMINI_DIR)..."
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
		rm -f "$$target"; \
		\
		gtarget=$(GEMINI_DIR)/$$d; \
		gtarget=$${gtarget%/}; \
		rm -f "$$gtarget"; \
	done
	@rm -f $(CLAUDE_DIR)/settings.json
	@echo "Removing broken symlinks..."
	@find $(CLAUDE_DIR)/agents $(CLAUDE_DIR)/rules $(CLAUDE_DIR)/commands $(CLAUDE_DIR)/skills -xtype l -delete 2>/dev/null || true
	@find $(GEMINI_DIR)/skills -xtype l -delete 2>/dev/null || true
	@echo "Done."
