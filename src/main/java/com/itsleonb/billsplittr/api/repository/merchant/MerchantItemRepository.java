package com.itsleonb.billsplittr.api.repository.merchant;

import com.itsleonb.billsplittr.api.entity.merchant.MerchantItem;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface MerchantItemRepository extends JpaRepository<MerchantItem, UUID> {
  boolean existsByMerchantIdAndName(UUID merchantId, String name);
}
