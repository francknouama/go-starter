# Homebrew Publishing Setup

This document explains how to set up the Personal Access Token (PAT) required for Homebrew tap publishing.

## The Issue

The release workflow needs to update the `francknouama/homebrew-tap` repository, but the default `GITHUB_TOKEN` only has permissions for the current repository (`francknouama/go-starter`).

## Solution: Personal Access Token (PAT)

### Step 1: Create a Personal Access Token

1. Go to GitHub → **Settings** → **Developer settings** → **Personal access tokens** → **Tokens (classic)**
2. Click **"Generate new token (classic)"**
3. Fill in the details:
   - **Note**: `go-starter-homebrew-publishing`
   - **Expiration**: Set to your preference (recommend 1 year)
   - **Scopes**: Select the following:
     - ✅ `repo` (Full control of private repositories)
     - ✅ `public_repo` (Access public repositories)

4. Click **"Generate token"**
5. **Important**: Copy the token immediately - you won't be able to see it again!

### Step 2: Add PAT as Repository Secret

1. Go to the `francknouama/go-starter` repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **"New repository secret"**
4. Enter:
   - **Name**: `HOMEBREW_TOKEN`
   - **Secret**: Paste your PAT from Step 1
5. Click **"Add secret"**

### Step 3: Verify Setup

After adding the secret, the next release should automatically publish to the Homebrew tap without the 403 permission error.

## Files Modified

The following files have been updated to support PAT-based authentication:

- `.github/workflows/release.yml` - Added `HOMEBREW_TOKEN` environment variable
- `.goreleaser.yml` - Added `token: "{{ .Env.HOMEBREW_TOKEN }}"` to brews configuration

## Testing

To test the setup:
1. Create a test tag: `git tag v1.3.2-test && git push origin v1.3.2-test`
2. Check the release workflow logs
3. Verify the homebrew formula is updated in `francknouama/homebrew-tap`
4. Delete the test tag: `git tag -d v1.3.2-test && git push origin :refs/tags/v1.3.2-test`

## Troubleshooting

If you still see 403 errors:
- Verify the PAT has `repo` scope
- Ensure the secret name is exactly `HOMEBREW_TOKEN`
- Check that the PAT hasn't expired
- Confirm you have admin access to the `homebrew-tap` repository