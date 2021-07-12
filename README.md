# Update Slack Info - GitHub Action
[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)


A [GitHub Action](https://github.com/features/actions) to update the Slack Information.

## Usage

You can use this action after any other action. Here is an example setup of this action:

1. Create a `.github/workflows/update-slack-info.yml` file in your GitHub repo.
2. Add the following code to the `update-slack-info.yml` file.

```yml
on: push
name: Update Slack Info Demo
jobs:
  slackNotification:
    name: Update Slack Info
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Update Slack Info
      uses: giantswarm/update-slack-info@v2
      env:
        SLACK_TOKEN: ${{ secrets.SLACK_TOKEN }}
```

3. Create `SLACK_TOKEN` secret using [GitHub Action's Secret](https://help.github.com/en/actions/configuring-and-managing-workflows/creating-and-storing-encrypted-secrets#creating-encrypted-secrets-for-a-repository). You can [generate a Slack token from here](https://slack.com/intl/en-gb/help/articles/215770388-Create-and-regenerate-API-tokens).
