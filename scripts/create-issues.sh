#!/bin/bash

# Ensure gh CLI is installed and authenticated
if ! command -v gh &> /dev/null; then
    echo "GitHub CLI (gh) is not installed. Please install it to proceed."
    exit 1
fi

echo "Creating backend issues..."

gh issue create \
  --title "implement backend email templates data layer and validation schemas" \
  --body "**Description:** Add structural data models and DTO validation rules for managing email templates.

**Checklist:**
* [ ] Create \`internal/models/template_dto.go\` defining CreateTemplateRequest and TemplateResponse with validation tags.
* [ ] Implement \`internal/repositories/template_repo.go\` wrapping PocketBase reads/writes using schema constants.
* [ ] Create \`internal/services/template_service.go\` handling domain rules for template management." \
  --label "enhancement"

gh issue create \
  --title "implement email templates HTTP handlers and routing endpoints" \
  --body "**Description:** Map HTTP requests to the email templates domain service layer and wire up authenticated routing.

**Checklist:**
* [ ] Create \`internal/handlers/templates_handler.go\` to bind requests and map typed errors to HTTP responses.
* [ ] Register routes in \`internal/handlers/routes.go\` applying the authentication middleware layer.
* [ ] Write unit tests for TemplateService utilizing fakes/mocks." \
  --label "enhancement"

gh issue create \
  --title "implement backend campaigns data layer and validation schemas" \
  --body "**Description:** Add structural data models and DTO validation rules for managing email blasts and campaigns.

**Checklist:**
* [ ] Create \`internal/models/campaign_dto.go\` defining CampaignRequest and CampaignResponse.
* [ ] Implement \`internal/repositories/campaign_repo.go\` for campaign persistence.
* [ ] Create \`internal/services/campaign_service.go\` handling business logic like configurable batching and scheduling." \
  --label "enhancement"

gh issue create \
  --title "implement campaigns HTTP handlers and routing endpoints" \
  --body "**Description:** Map HTTP requests to the campaigns domain service layer to trigger and manage email blasts.

**Checklist:**
* [ ] Create \`internal/handlers/campaigns_handler.go\` to handle campaign execution endpoints.
* [ ] Register routes in \`internal/handlers/routes.go\` with appropriate authentication.
* [ ] Write unit tests for CampaignService." \
  --label "enhancement"

gh issue create \
  --title "implement backend sender accounts data layer and validation schemas" \
  --body "**Description:** Add data models and integration logic for managing Google Auth sender accounts.

**Checklist:**
* [ ] Create \`internal/models/account_dto.go\` defining AccountRequest and AccountResponse.
* [ ] Implement \`internal/repositories/account_repo.go\` for sender account persistence.
* [ ] Create \`internal/services/account_service.go\` to handle OAuth token management and sending constraints." \
  --label "enhancement"

gh issue create \
  --title "implement sender accounts HTTP handlers and routing endpoints" \
  --body "**Description:** Map HTTP requests to the sender accounts domain service layer.

**Checklist:**
* [ ] Create \`internal/handlers/accounts_handler.go\` for OAuth flows and account management.
* [ ] Register routes in \`internal/handlers/routes.go\`.
* [ ] Write unit tests for AccountService." \
  --label "enhancement"

echo "Creating frontend issues..."

gh issue create \
  --title "implement connections dashboard UI and data fetching" \
  --body "**Description:** Build the CRM-style directory for adding and managing network connections using React Router v7 and React Query.

**Checklist:**
* [ ] Create \`web/app/routes/connections.tsx\` defining the route and layout.
* [ ] Implement data fetching and mutations in the route using React Query Kit.
* [ ] Build the connection list and add/edit forms using Tailwind CSS and Framer Motion." \
  --label "enhancement"

gh issue create \
  --title "implement email templates builder UI and validation forms" \
  --body "**Description:** Build the interface for creating and managing reusable email templates.

**Checklist:**
* [ ] Create \`web/app/routes/templates.tsx\` and nested builder routes.
* [ ] Implement form validation using standard React Router form actions/React Query mutations.
* [ ] Build the template preview and editor components." \
  --label "enhancement"

gh issue create \
  --title "implement email blast campaign execution UI" \
  --body "**Description:** Build the execution engine UI where users select templates, choose recipients, and configure batching settings.

**Checklist:**
* [ ] Create \`web/app/routes/campaigns.tsx\` defining the campaign launch flow.
* [ ] Implement the batching configuration UI and recipient selection steps.
* [ ] Connect to backend campaign execution endpoints." \
  --label "enhancement"

gh issue create \
  --title "implement sender accounts integration UI" \
  --body "**Description:** Build the settings interface for users to connect and manage their Google sender accounts.

**Checklist:**
* [ ] Create \`web/app/routes/accounts.tsx\` for sender profile management.
* [ ] Implement Google OAuth connection flows in the UI." \
  --label "enhancement"

gh issue create \
  --title "implement application analytics dashboard UI" \
  --body "**Description:** Build the high-level overview dashboard showing email performance, open rates, and delivery success.

**Checklist:**
* [ ] Create \`web/app/routes/analytics.tsx\` for the main dashboard view.
* [ ] Implement data visualizations and metrics cards leveraging Tailwind CSS and Framer Motion." \
  --label "enhancement"

echo "All issues created successfully!"
