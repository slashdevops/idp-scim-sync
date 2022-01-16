# Using SSO

**Warning:** Comments here are [IMHO](https://en.wiktionary.org/wiki/IMHO)

Start using [Single Sign-On (SSO)](https://en.wikipedia.org/wiki/Single_sign-on) in your AWS Accounts is hard to achieve, you must need to read a lot of documentation and of course, you need to know how to use the [AWS SSO service](https://aws.amazon.com/blogs/security/how-to-create-and-manage-users-within-aws-sso/).

>For me, __technical things__ are just things that are going to work sooner or later or as soon you understand how to use these.  But __non-technical things__ are just things that most of the time you will improve according to __your time__ implementing __technical things__ and this could we call __"experience"__.

So, said that, let me help you with my `"experience"` using SSO.

## Recommendations

Before start

* Do a planning of `how many services` you will need to integrate with your `Identity Provider using SSO`, and `how many users and Groups you will need to create`.
* Use `Groups and their Members` to create `Users` in your SSO integration, `avoid the integration directly with Users`.
* Establish a `Naming Convention` for your `SSO Groups`.
* Use prefixes for your `SSO Groups name`.

A little bit more

* Google Workspace Groups and AWS Single Sign-On are free of charge, so take advantage of this
* Filter the data in the `source` is always better than in the `process`
* [WYSIWYG](https://dictionary.cambridge.org/es/diccionario/ingles/wysiwyg) is better than `opaque or shadowed process`, I mean, if you have 1 group called `My Group` with `2 members`, `user.1@mydomain.com` and `user.2@mydomain.com` in the `Google Workspace`, everybody is expecting to see the same in `AWS SSO side` regardless of what this program does
* A process that scales independently of the `source` is better than a `process that needs too many changes during its escalation`

### Example

Given Google Workspace Groups with these conditions:

| Group Name         | Group Email                    | Members                  |
| ------------------ | ------------------------------ | ------------------------ |
| AWS Administrators | aws-administrator@mydomain.com | [a,b,c,d]@mydomain.com   |
| AWS DevOps         | aws-devops@mydomain.com        | [f,g,h]@mydomain.com     |
| AWS Developers     | aws-developers@mydomain.com    | [k,j,z,a,g]@mydomain.com |
| ...                | ...                            | ...                      |

your [idpscim](https://github.com/slashdevops/idp-scim-sync/blob/main/docs/idpscim.md) `[AWS Lambda function|container image|cli]` could use `--gws-groups-filter 'name=AWS* email:aws-*'`

This is easy and it is in compliance with the previous recommendations, but I think the most important ones are:
> You can `increase or decrease` the number of `groups and their members` in __Google Workspace__ and never need to `change the parameters` of the __idpscim__ `[AWS Lambda function|container image|cli]`

## TL;DR

__NOTES:__

* This is a `WIP`, keep calm and don't panic
* [TL;DR means: Too Long Didn't Read](https://en.wikipedia.org/wiki/TL;DR)

### Planning

### Groups and members

### Use Naming Conventions

The most important part of implementing SSO is the planning of the things you need to do before that, but for me the most important one is `Groups name and Naming convention`,

__Why?__

Because one you `sync the groups and users` the first time with [AWS SSO](https://aws.amazon.com/single-sign-on) and assign to them the [Permission sets](https://docs.aws.amazon.com/singlesignon/latest/userguide/permissionsetsconcept.html) you should not change the `Groups and Users main attributes`, with in case of this program are:

| Entity | Main Attribute |
| ------ | -------------- |
| Groups | Name           |
| Users  | Email          |

Of course, you can change the `Groups and Users main attributes` but after that you will lose the `Permission sets` assigned to them, so __you should not change them__.

__Again, Why?__

>Imagine you have __40 Groups__ and __1,000 Users__ and some of these __Groups__ has __50 Users__ and you have assigned some [Permission sets](https://docs.aws.amazon.com/singlesignon/latest/userguide/permissionsetsconcept.html) to this `Group`, so if you __change the Name__ of __this Group__ in your Identity Provider, our case __Google Workspace Directory Service__, __50 users of your company will be lost their permissions__ to your __AWS Accounts in just one moment__, and this could be __hard to fix__ is you don't have implemented any [IaC](https://en.wikipedia.org/wiki/Infrastructure_as_code) practice to match the __Groups__ and __Users__ with the __Permission sets__ assigned to them.

So, to avoid this scenario you should: __Plan your Groups Naming Convention First__

### Use Name Prefixes
