package com.itsleonb.billsplittr.impl.controller;

import com.itsleonb.billsplittr.api.controller.AuthController;
import com.itsleonb.billsplittr.api.model.JsonResponse;
import com.itsleonb.billsplittr.api.model.auth.LoginRequest;
import com.itsleonb.billsplittr.api.model.auth.LoginResponse;
import com.itsleonb.billsplittr.api.model.auth.RegisterRequest;
import com.itsleonb.billsplittr.impl.service.auth.AuthServiceImpl;
import lombok.AllArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RestController;

@RestController
@AllArgsConstructor(onConstructor = @__(@Autowired))
public class AuthControllerImpl implements AuthController {
  private AuthServiceImpl authService;

  @Override
  public JsonResponse<String> handleRegister(RegisterRequest request) {
    authService.register(request);

    return JsonResponse.<String>builder()
      .success(true)
      .data("Register success, please login")
      .build();
  }

  @Override
  public JsonResponse<LoginResponse> handleLogin(LoginRequest request) {
    LoginResponse response = authService.login(request);

    return JsonResponse.<LoginResponse>builder()
      .success(true)
      .data(response)
      .build();
  }
}
