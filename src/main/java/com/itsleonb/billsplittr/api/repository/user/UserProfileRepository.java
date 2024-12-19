package com.itsleonb.billsplittr.api.repository.user;

import com.itsleonb.billsplittr.api.entity.user.UserProfile;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface UserProfileRepository extends JpaRepository<UserProfile, String> {
  Optional<UserProfile> findByUserId(UUID userId);
}
