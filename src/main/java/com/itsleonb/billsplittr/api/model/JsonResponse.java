package com.itsleonb.billsplittr.api.model;

import com.fasterxml.jackson.annotation.JsonInclude;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@JsonInclude(JsonInclude.Include.NON_NULL)
public class JsonResponse<T> {
  private boolean success;
  private T data;
  private ErrorResponse error;

  public static JsonResponse<String> NewErrorResponse(Exception e) {
    return JsonResponse.<String>builder()
      .success(false)
      .error(ErrorResponse.builder()
        .type(e.getClass().getSimpleName())
        .message(e.getMessage())
        .build())
      .build();
  }
}
