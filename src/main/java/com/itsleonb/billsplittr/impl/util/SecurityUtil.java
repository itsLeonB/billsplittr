package com.itsleonb.billsplittr.impl.util;

import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Component;

@Component
public class SecurityUtil {
  public String getCurrentUserEmail() {
    Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
    if (authentication == null || !authentication.isAuthenticated()) {
      throw new SecurityException("Unauthenticated");
    }

    UserDetails userDetails = (UserDetails) authentication.getPrincipal();
    
    return userDetails.getUsername();
  }
}
