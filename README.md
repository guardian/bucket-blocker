# bucket-blocker

## What is this thing?

Bucket blocker is a tool that searces for violations of [FSBP S3.8](https://docs.aws.amazon.com/securityhub/latest/userguide/s3-controls.html#s3-8), which states that all buckets should have individual configurations blocking public access.

Once it's found the buckets that are not compliant, it skips over any that are in CloudFormation stacks (to avoid introducing stack drift), and then blocks public access to the remaining buckets.

```mermaid
flowchart TB
    stack[Is it part of a cloudformation stack]
    excl[Is it in a list of excluded \n buckets provided by the user?]
    block[Block public access to the bucket]
    ruleBreak[Does the bucket break S3.8?]
    break[Do nothing.]
    noAccess[No. Access already \n blocked]

    ruleBreak --> Yes --> stack --> No --> excl --> Nope --> block
    ruleBreak --> noAccess --> break
    stack --> Yeah --> break
    excl --> Yep --> break
```

There are a few extra features, controlled by flags, enumerated below.

## Command line options:

Bucket blocker takes up to 3 flags:

- **profile**: _Required._ The profile to use when connecting to AWS.

- **region**: _Required._ The region where the bucket is located.

- **execute**: _Optional._ Takes no value. If present, it will ask the user to confirm, then block the buckets. If not, it will only print
  the buckets that would have been blocked.

- **exclusions**: _Optional._ Comma-delimited list of buckets to exclude from blocking.

- **max**: _Optional._ The maximum number of buckets to block. Between 1
  and 100. Defaults to 100, which is the maximum number of buckets that can
  exist in an AWS account.

You will also need credentials for the relevant AWS account from Janus.

## Running the binary

This application is downloadable from brew. You'll need the guardian's brew tap
installed before you can install the application.

To do this all at once, run the following command:

```bash
brew tap guardian/homebrew-devtools && brew install bucket-blocker
```

You can also download the binary directly from the Releases page on GitHub, or
build it from source.

## Local development

<!-- TODO enforce conventional commits via GHA-->

When committing your changes, please use the
[conventional commit](https://www.conventionalcommits.org/en/v1.0.0/#summary)
format. This will allow us to automatically generate a changelog and correctly
version the application when it is released.

While developing locally, you can test the application using the following
command from the root of the repository, without needing to build the binary:

```bash
go run main.go -profile <PROFILE> -region <REGION> [OPTIONAL_FLAGS]
```

## Releasing to brew

Creating a new release of the application on brew, is currently a manual
process. You will need to update the version, urls, and SHAs in
[this file](https://github.com/guardian/homebrew-devtools/blob/main/Formula/bucket-blocker.rb)
in the homebrew-devtools repo.
