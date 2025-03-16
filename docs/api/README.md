# Badge Assignment System API Documentation

This directory contains complete documentation for the Badge Assignment System API.

## Contents

- [Overview](#overview)
- [Authentication](#authentication)
- [Public API](#public-api)
- [Admin API](#admin-api)
- [Examples](#examples)
- [Error Handling](#error-handling)

## Overview

The Badge Assignment System API is a RESTful API that allows clients to manage badges, events, and user-badge assignments. The API is divided into two main sections:

1. **Public API Endpoints**: These endpoints are accessible to regular users and provide read access to badges and functionality for submitting events.
2. **Admin API Endpoints**: These endpoints are for administrative operations and require appropriate authentication.

## API Reference

For detailed API documentation, please refer to the following sections:

- [Badge API Documentation](./badges.md) - Badge management endpoints
- [Event API Documentation](./events.md) - Event handling endpoints
- [User Badge API Documentation](./user-badges.md) - User-badge relationship endpoints
- [Event Type API Documentation](./event-types.md) - Event type management endpoints
- [Condition Type API Documentation](./condition-types.md) - Condition type management endpoints

## Authentication

Currently, the API does not enforce authentication. Future versions will include proper authentication and authorization mechanisms.

## Examples

For examples of API usage, see the [examples directory](./examples/). 
