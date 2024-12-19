package com.itsleonb.billsplittr.api.service;

import com.itsleonb.billsplittr.api.model.auth.LoginRequest;
import com.itsleonb.billsplittr.api.model.auth.LoginResponse;
import com.itsleonb.billsplittr.api.model.auth.RegisterRequest;

public interface AuthService {
  void register(RegisterRequest request);

  LoginResponse login(LoginRequest request);
}
