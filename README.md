# git-copy

## Usage

```yaml
name: Copy Setting Deploy - Dev
on:
  push:
    branches:
      - master
    paths:
      - "settings/**"
  workflow_dispatch:

concurrency:
  group: dev-setting-deploy-${{ github.ref }}

jobs:
  deploy:
    runs-on: ubuntu-latest
    permissions:
      contents: "read"
      id-token: "write"
    env:
      LOCATION: europe-west1

    steps:
      - name: Checkout repo
        uses: actions/checkout@v3

      - name: Copy to repo
        uses: pal-paul/git-copy@v1
        with:
          owner: "group-digital"
          repo: "destination-repo"
          token: "${{ secrets.GH_TOKEN }}"
          branch: "story/update-settings"
          pull_message: "update settings"
          file_path: "./settings/markets.json"
          detination_file_path: "configs/dev.json"
          team_reviewers: "reviewer-team-name"
```