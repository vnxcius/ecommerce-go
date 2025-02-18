package main

type contextKey string

const isAuthenticatedUserKey = contextKey("isAuthenticatedUser")
const userIDKey = "userID"