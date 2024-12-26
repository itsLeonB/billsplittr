package com.itsleonb.billsplittr.api.repository.merchant;

import com.itsleonb.billsplittr.api.entity.merchant.Merchant;
import org.springframework.data.jpa.repository.EntityGraph;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface MerchantRepository extends JpaRepository<Merchant, UUID> {
  @Query("SELECT m FROM Merchant m WHERE LOWER(m.name) LIKE LOWER(CONCAT(:name, '%'))")
  List<Merchant> searchByName(@Param("name") String name);

  boolean existsByName(String name);

  @EntityGraph(attributePaths = {"items"})
  Optional<Merchant> findWithItemsById(UUID id);
}
