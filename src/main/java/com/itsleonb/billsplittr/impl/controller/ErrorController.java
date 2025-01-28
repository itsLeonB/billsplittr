package com.itsleonb.billsplittr.impl.controller;

import com.itsleonb.billsplittr.api.exception.ConflictException;
import com.itsleonb.billsplittr.api.exception.NotFoundException;
import com.itsleonb.billsplittr.api.model.JsonResponse;
import jakarta.validation.ConstraintViolationException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.http.converter.HttpMessageNotReadableException;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.security.core.AuthenticationException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@RestControllerAdvice
public class ErrorController {
  @ExceptionHandler(ConstraintViolationException.class)
  public ResponseEntity<JsonResponse<String>> constraintViolationException(ConstraintViolationException exception) {
    return errorResponse(exception, HttpStatus.BAD_REQUEST);
  }

  @ExceptionHandler(ConflictException.class)
  public ResponseEntity<JsonResponse<String>> conflictException(ConflictException exception) {
    return errorResponse(exception, HttpStatus.CONFLICT);
  }

  @ExceptionHandler(NotFoundException.class)
  public ResponseEntity<JsonResponse<String>> notFoundException(NotFoundException exception) {
    return errorResponse(exception, HttpStatus.NOT_FOUND);
  }

  @ExceptionHandler(AuthenticationException.class)
  public ResponseEntity<JsonResponse<String>> authenticationException(AuthenticationException exception) {
    return errorResponse(exception, HttpStatus.UNAUTHORIZED);
  }

  @ExceptionHandler(AccessDeniedException.class)
  public ResponseEntity<JsonResponse<String>> accessDeniedException(AccessDeniedException exception) {
    return errorResponse(exception, HttpStatus.FORBIDDEN);
  }

  @ExceptionHandler(HttpMessageNotReadableException.class)
  public ResponseEntity<JsonResponse<String>> handleMessageNotReadable(HttpMessageNotReadableException ex) {
    return errorResponse(new RuntimeException("Invalid request body"), HttpStatus.BAD_REQUEST);
  }

  private ResponseEntity<JsonResponse<String>> errorResponse(Exception e, HttpStatus code) {
    return ResponseEntity.status(code).body(JsonResponse.newErrorResponse(e));
  }
}
