package com.itsleonb.billsplittr.impl.controller;

import com.itsleonb.billsplittr.api.exception.BadRequestException;
import com.itsleonb.billsplittr.api.exception.ConflictException;
import com.itsleonb.billsplittr.api.exception.NotFoundException;
import com.itsleonb.billsplittr.api.model.JsonResponse;
import io.jsonwebtoken.ExpiredJwtException;
import jakarta.validation.ConstraintViolationException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.http.converter.HttpMessageNotReadableException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.sql.SQLException;

@Slf4j
@RestControllerAdvice
public class ErrorController {
  @ExceptionHandler({
    ConstraintViolationException.class,
    HttpMessageNotReadableException.class,
    BadRequestException.class
  })
  public ResponseEntity<JsonResponse<String>> badRequestException(Exception exception) {
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

  @ExceptionHandler(ExpiredJwtException.class)
  public ResponseEntity<JsonResponse<String>> unauthorizedException(Exception exception) {
    return errorResponse(exception, HttpStatus.UNAUTHORIZED);
  }

  @ExceptionHandler({SQLException.class, Exception.class})
  public ResponseEntity<JsonResponse<String>> internalServerException(Exception exception) {
    for (StackTraceElement stackTraceElement : exception.getStackTrace()) {
      log.info(stackTraceElement.toString());
    }

    return errorResponse(exception, HttpStatus.INTERNAL_SERVER_ERROR);
  }

  private ResponseEntity<JsonResponse<String>> errorResponse(Exception e, HttpStatus code) {
    return ResponseEntity.status(code).body(JsonResponse.NewErrorResponse(e));
  }
}
