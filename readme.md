# Knowledge trees and courses

This repo holds the information about knowledge trees for all topics and course contents.

Course versioning is done through commit SHA specification. The new entry is created in DB for model CourseVersion where repo, path and SHA are specified for the version of the course to be used.

Courses are also differentiated by software versions used. Separate files are created for each new software version specified.

```
/backend/intro_python/
                        - python_v1.yaml #holds course for python x and click x versions
                        - python_v2.yaml #holds course for python x+1 and click x+1 versions
```

## Course creation workflow

1. Course outline created by using nodes from knowledge domain tree
2. Course object and CourseVersion are created in DB specifying repo, owner, path and SHA for concrete course version
3. Background task is scheduled to check the status of a course lessons (if they are present, correct version is there and they are published)
4. Task Result object is created in DB with specification of missing lessons and their general status
5. After all lessons were created the background task is scheduled again and course status changes to ready
6. When course status is ready it can be moved to "published" status

## Dependencies and lessons

Domain only contains information about node names and ids. Actual lessons, versions and software versions are created through API calls to lessons endpoint.
Software dependencies defined for one node are applied to all children nodes and their children recursively. When new lesson version is created through API calls the following checks are done:

- ID exists in knowledge domain tree
- The software versions defined match the node restrictions on software version
  After new lesson version was created and all checks passed, the lesson version goes to "ready" status and could "published" laster

## How it works?

- If SHA for course, sprint, lesson or block version gets changed, new ID is generated
