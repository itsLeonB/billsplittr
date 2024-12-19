package com.itsleonb.billsplittr.impl.service.auth;

import com.itsleonb.billsplittr.api.entity.user.User;
import com.itsleonb.billsplittr.api.entity.user.UserProfile;
import com.itsleonb.billsplittr.api.exception.ConflictException;
import com.itsleonb.billsplittr.api.exception.NotFoundException;
import com.itsleonb.billsplittr.api.model.auth.LoginRequest;
import com.itsleonb.billsplittr.api.model.auth.LoginResponse;
import com.itsleonb.billsplittr.api.model.auth.RegisterRequest;
import com.itsleonb.billsplittr.api.repository.user.UserProfileRepository;
import com.itsleonb.billsplittr.api.repository.user.UserRepository;
import com.itsleonb.billsplittr.api.service.AuthService;
import jakarta.validation.ConstraintViolation;
import jakarta.validation.ConstraintViolationException;
import jakarta.validation.Validator;
import lombok.AllArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Set;

@Service
@AllArgsConstructor(onConstructor_ = @__(@Autowired))
public class AuthServiceImpl implements AuthService {
  private final AuthenticationManager authenticationManager;
  private final UserProfileRepository userProfileRepository;
  private final PasswordEncoder passwordEncoder;
  private final UserRepository userRepository;
  private final JwtService jwtService;
  private final Validator validator;

  @Override
  @Transactional
  public void register(RegisterRequest request) {
    Set<ConstraintViolation<RegisterRequest>> violations = validator.validate(request);
    if (!violations.isEmpty()) {
      throw new ConstraintViolationException(violations);
    }

    if (userRepository.existsByEmail(request.getEmail())) {
      throw new ConflictException("Email already registered");
    }

    User user = User.builder()
      .email(request.getEmail())
      .password(passwordEncoder.encode(request.getPassword()))
      .build();

    User savedUser = userRepository.save(user);

    UserProfile profile = UserProfile.builder()
      .userId(savedUser.getId())
      .build();

    userProfileRepository.save(profile);
  }

  @Override
  public LoginResponse login(LoginRequest request) {
    Set<ConstraintViolation<LoginRequest>> violations = validator.validate(request);
    if (!violations.isEmpty()) {
      throw new ConstraintViolationException(violations);
    }

    authenticationManager.authenticate(
      new UsernamePasswordAuthenticationToken(
        request.getEmail(),
        request.getPassword()
      )
    );

    User user = userRepository
      .findByEmail(request.getEmail())
      .orElseThrow(() -> new NotFoundException("wrong email/password combination"));

    return LoginResponse.builder()
      .token(jwtService.generateToken(user))
      .build();
  }
}
