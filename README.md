# Probe

A small tracker for Github issues and PRs.
At the moment it does not automatically update, but that should be added at some point.

You will need a config file either in `$HOME/.probe.yml` with two keys

```yaml
github_token: your-github-token
database_path: /path/to/where/you/want/to/store/the_db.json
```

# Usage

From there you can list all current tracked items with `probe ls`:

```
IDX     TITLE                                                                                           STATUS  LAST ACTION             LAST CHANGED    LINK                                                     
18      Feature Request: Security Hub                                                                OPEN    CrossReferencedEvent    10 hours ago  https://github.com/hashicorp/terraform-provider-aws/issues/6674
19      r/aws_securityhub: Add aws_securityhub_invite_accepter resource                              OPEN    SubscribedEvent         1 week ago    https://github.com/hashicorp/terraform-provider-aws/pull/12684
20      AWS SecurityHub: the ability to disable specific compliance controls                         OPEN    LabeledEvent            2 months ago  https://github.com/hashicorp/terraform-provider-aws/issues/11624
21      aws_transfer_user: Support "restricted" home directory                                       CLOSED  LockedEvent             1 month ago   https://github.com/hashicorp/terraform-provider-aws/issues/11632
22      aws_transfer_user should provide home-directory-mappings option.                             CLOSED  LockedEvent             1 month ago   https://github.com/hashicorp/terraform-provider-aws/issues/11281
23      Transfer Server VPC Endpoint Type and User Home Directory Type / Mapping                     MERGED  LockedEvent             1 month ago   https://github.com/hashicorp/terraform-provider-aws/pull/12599
24      resource/aws_transfer_user: add home_directory_type and home_directory_mappings arguments    MERGED  LockedEvent             1 month ago   https://github.com/hashicorp/terraform-provider-aws/pull/13591
```

Adding new items with `probe add https://github.com/hashicorp/terraform-provider-aws/pull/123` or a similar url for links.

If you want to delete an item, you can call `probe del 18` to delte the item with `IDX` 18.

Lastly, regularly call `probe update` to see if the status of the PR/Issue has changed.
