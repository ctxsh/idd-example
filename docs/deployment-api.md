# Deployment API

## Summary

Build a small HTTP API that lets a team request, approve, and track software
deployments. The API is a lightweight control plane: it records deployment intent
and lifecycle state, but it does not actually deploy infrastructure.

## Problem

Teams often coordinate deployments across issues, pull requests, chat, and
manual checklists. That makes it hard to answer basic questions:

- What is being deployed?
- Who requested it?
- Has it been approved?
- What risk was accepted?
- What is the rollback plan?
- What happened after the deployment started?

This project should make those answers explicit and reviewable.

## Target Users

- Engineers requesting a deployment.
- Reviewers approving or rejecting a deployment.
- On-call engineers checking current deployment state.

## Initial Scope

The API should support:

- Creating a deployment request.
- Listing deployment requests.
- Fetching one deployment request.
- Approving or rejecting a deployment request.
- Moving an approved deployment through started, succeeded, failed, and rolled
  back states.
- Recording issue or pull request references on the deployment.
- Recording risk and rollback notes before production deployment.

## Out Of Scope

- Running real deploys.
- Integrating with Kubernetes, cloud APIs, or CI systems.
- Full authentication and authorization.
- Web UI.
- Multi-service orchestration.

## Deployment Model

A deployment request should include:

- `id`
- `service`
- `environment`
- `version`
- `status`
- `requested_by`
- `approved_by`
- `issue`
- `risk`
- `rollback_plan`
- `created_at`
- `updated_at`

Possible statuses:

- `pending_approval`
- `approved`
- `rejected`
- `started`
- `succeeded`
- `failed`
- `rolled_back`

## Candidate API

- `POST /deployments`
- `GET /deployments`
- `GET /deployments/{id}`
- `POST /deployments/{id}/approve`
- `POST /deployments/{id}/reject`
- `POST /deployments/{id}/start`
- `POST /deployments/{id}/succeed`
- `POST /deployments/{id}/fail`
- `POST /deployments/{id}/rollback`

## Design Constraints

- Prefer clear lifecycle rules over infrastructure realism.
- Make invalid state transitions explicit errors.
- Make issue references visible in API responses so PRs can connect back to the
  planned work.
