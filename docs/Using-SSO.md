# Using SSO

**Warning:** Comments here are [IMHO](https://en.wiktionary.org/wiki/IMHO)

Start using [Single Sign-On (SSO)](https://en.wikipedia.org/wiki/Single_sign-on) in your AWS Accounts is hard to achieve, you must need to read a lot of documentation and of course, you need to know how to use the [AWS SSO service](https://aws.amazon.com/blogs/security/how-to-create-and-manage-users-within-aws-sso/).

>For me, __technical things__ are just things that are going to work sooner or later or as soon you understand how to use these.  But __non-technical things__ are just things that most of the time you will improve according to __your time__ implementing __technical things__ and this could we call __"experience"__.

So, said that, let me help you with my "experience" using SSO.

## Recommendations

Before start

* Do a planning of `how many services` you will need to integrate with your `Identity Provider using SSO`, and `how many users and Groups you will need to create`.
* Use `Groups and their members` to create Users in your SSO integration, `avoid the integration directly with users`.
* Establish a `Naming Convention` for your `SSO Groups`.
* Use prefixes for your SSO Groups name.

## TL;DR

__NOTE:__ [TL;DR means: Too Long Didn't Read](https://en.wikipedia.org/wiki/TL;DR)

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
