# Konex: Keeper of Network Email Execution

**Konex** (Keeper of Network Email Execution) is an automated email marketing and personal networking tool designed to help individuals maintain connections, showcase their work, and execute targeted email campaigns with ease.

At its core, Konex empowers users to build and nurture their professional network by automating the delivery of customized email templates to leads and reachouts, without risking their domain reputation.

---

## 🎯 The Motive

Building and maintaining a professional network can be time-consuming. Whether you're reaching out to your connections to ask for a referral, sharing your latest projects, or circulating a recent LinkedIn post, doing it manually is tedious.

Konex solves this by allowing you to:

- Curate your network (Leads/Reachouts).
- Create reusable email templates.
- Automatically execute "Email Blasts" to your selected network with a single click.
- Send monthly updates to keep your network engaged.

---

## ✨ Key Features & Capabilities

### 1. Automated Email Execution

Create an email template (e.g., a "Referral Request") and send it to your entire network with just one click. Konex handles the automated delivery in the background.

### 2. Network Nurturing & Showcasing

Easily showcase your latest projects, blog articles, or LinkedIn posts by scheduling monthly reminder emails to your contacts, keeping you at the top of their minds.

### 3. Smart Configurable Batching

To protect your email domain's sender reputation and avoid being flagged as spam, Konex features **Configurable Batching**.
_Example:_ You can configure the system to send only 1 email every 10 minutes, spreading out a blast of 50 emails over several hours.

### 4. Seamless Google Authentication

Connect your sender accounts securely using Google Auth, allowing Konex to send emails directly on your behalf.

---

## 🧩 Modules & Architecture

The application is structured around a centralized dashboard offering the following modules:

- **📊 Analytics:** A high-level overview of your email performance, open rates, delivery success, and network engagement.
- **👥 Connections:** A CRM-style directory where you can add, organize, and manage the people in your network.
- **📝 Email Templates:** A module to create, edit, and store reusable email templates for different scenarios (referrals, showcases, updates).
- **🚀 Email Blast:** The execution engine where you select a template, choose your recipients, configure your batching settings, and launch the campaign.
- **📧 Sender Accounts:** Manage the email addresses you use to send campaigns, integrated securely via Google Auth.
- **⚙️ Settings:** Global application settings, user preferences, and default configuration options.
