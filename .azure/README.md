# Azure DevOps

Azure DevOps is just a reincarnation of Visual Studio Team Services...
laugh all you want :)

The reason we have moved to Azure DevOps (and away from CircleCI) is due to the
fact that it offers free Linux, MacOs & Windows machines.

Yes TravisCI has also recently added Windows support and is the only other CI
provider that I know of that offers all 3 major operating systems. However in
my opinion Azure's offering is better than that offered by Travis.

> UPDATE (August 2019): CircleCI now offer Windows based virtual machines
> but that is not offered in any free plan.

## Folder Layout

- `./.azure/executors`:
  Is for any Dockerfiles that might be used for container based jobs.
  The docker experience on Azure DevOps is sub optimal but useable if needed.

- `./.azure/pipelines`:
  Is where we keep all root pipeline definitions, yes Azure DevOps lets us
  define multiple pipelines for a single repo and trigger the pipelines based
  on what actually changed in any given commit. I think this a great feature.
  Although not required for this project at the moment.

- `./.azure/templates/steps`:
  Step templates live in here, see https://docs.microsoft.com/en-us/azure/devops/pipelines/process/templates?view=vsts#step-re-use

- `./.azure/templates/jobs`:
  Job templates live in here, see https://docs.microsoft.com/en-us/azure/devops/pipelines/process/templates?view=vsts#job-reuse
