package com.itsleonb.billsplittr.api.repository.user;

import com.itsleonb.billsplittr.api.entity.user.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface UserRepository extends JpaRepository<User, String> {
  boolean existsByEmail(String email);

  Optional<User> findByEmail(String email);
}
