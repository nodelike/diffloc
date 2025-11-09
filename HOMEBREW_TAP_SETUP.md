# Homebrew Tap Setup Guide

This guide explains how to set up the Homebrew tap for `diffloc` to enable automated releases.

## Prerequisites

1. GitHub account with access to create repositories
2. GoReleaser installed locally (for testing): `brew install goreleaser`
3. Personal Access Token (PAT) with `repo` and `workflow` permissions

## Step 1: Create the Homebrew Tap Repository

1. Go to GitHub and create a new repository named `homebrew-tap`
   - Repository name MUST be: `homebrew-tap`
   - Username/Org: `nodelike` (as specified in `.goreleaser.yml`)
   - Full repo path: `github.com/nodelike/homebrew-tap`
   - Make it public
   - Initialize with a README

2. Clone the repository locally:
   ```bash
   git clone https://github.com/nodelike/homebrew-tap.git
   cd homebrew-tap
   ```

3. Create the Formula directory:
   ```bash
   mkdir -p Formula
   git add Formula
   git commit -m "Initialize Formula directory"
   git push origin main
   ```

## Step 2: Configure GitHub Secrets

1. Go to the `diffloc` repository settings
2. Navigate to Settings → Secrets and variables → Actions
3. Add a new repository secret:
   - Name: `HOMEBREW_TAP_GITHUB_TOKEN`
   - Value: Your GitHub Personal Access Token (PAT)
   
   To create a PAT:
   - Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Generate new token with these scopes:
     - `repo` (all)
     - `workflow`
   - Copy the token and save it as the secret

## Step 3: Test Locally (Optional but Recommended)

Before creating a real release, test GoReleaser locally:

```bash
cd /path/to/diffloc
goreleaser release --snapshot --clean
```

This will:
- Build binaries for all platforms
- Create archives
- Generate checksums
- But NOT publish anything (due to `--snapshot`)

Check the `dist/` directory for output.

## Step 4: Create Your First Release

1. Ensure all changes are committed:
   ```bash
   git add -A
   git commit -m "Prepare for v0.1.0 release"
   git push origin main
   ```

2. Create and push a version tag:
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   git push origin v0.1.0
   ```

3. GitHub Actions will automatically:
   - Trigger the release workflow (`.github/workflows/release.yml`)
   - Run GoReleaser
   - Build binaries for all platforms
   - Create a GitHub release with binaries
   - Update the Homebrew tap formula automatically

4. Check the Actions tab in your GitHub repository to monitor progress.

## Step 5: Verify Homebrew Installation

After the release workflow completes:

1. Check the `homebrew-tap` repository:
   - A new file should appear at `Formula/diffloc.rb`
   - It should contain the formula with correct URLs and SHA256 checksums

2. Test installation locally:
   ```bash
   brew tap nodelike/tap
   brew install nodelike/tap/diffloc
   ```

3. Verify it works:
   ```bash
   diffloc --help
   ```

4. Test in a project:
   ```bash
   cd /path/to/some/project
   diffloc
   ```

## Troubleshooting

### Issue: GoReleaser fails with "token validation error"

**Solution**: Ensure the `HOMEBREW_TAP_GITHUB_TOKEN` secret is set correctly with proper permissions.

### Issue: Formula not updated in tap repository

**Solution**: 
- Check that the tap repository is public
- Verify the repository name is exactly `homebrew-tap`
- Check the GoReleaser logs in GitHub Actions

### Issue: Brew install fails with 404

**Solution**: 
- Ensure the GitHub release was created successfully
- Verify the URLs in the formula are correct
- Check that the binaries were uploaded to the release

### Issue: SHA256 mismatch

**Solution**: This usually means the formula references an old SHA256. Wait for GoReleaser to update the formula automatically, or manually fix it in the tap repository.

## Updating the Formula

For subsequent releases, simply:

1. Make your changes
2. Commit and push
3. Create a new tag (e.g., `v0.2.0`)
4. Push the tag

GoReleaser will automatically update the formula in your tap.

## Manual Formula Updates

If you need to manually update the formula:

1. Clone the tap repository
2. Edit `Formula/diffloc.rb`
3. Update version, URL, and SHA256
4. Commit and push

To calculate SHA256 of a tarball:
```bash
sha256sum diffloc_0.1.0_darwin_arm64.tar.gz
```

## Submitting to Homebrew Core

Once your project matures (75+ stars, 30+ days old, stable), you can submit to homebrew-core:

1. Read the [Homebrew documentation](https://docs.brew.sh/Adding-Software-to-Homebrew)
2. Test your formula thoroughly
3. Run `brew audit --new-formula diffloc`
4. Submit a PR to [homebrew-core](https://github.com/Homebrew/homebrew-core)

Benefits of homebrew-core:
- Users can install with just `brew install diffloc` (no tap needed)
- More visibility
- More downloads
- Part of the official Homebrew package index

## Resources

- [GoReleaser Documentation](https://goreleaser.com/intro/)
- [Homebrew Tap Documentation](https://docs.brew.sh/Taps)
- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)

