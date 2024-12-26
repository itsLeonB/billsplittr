package com.itsleonb.billsplittr.impl.service;

import com.itsleonb.billsplittr.api.entity.merchant.Merchant;
import com.itsleonb.billsplittr.api.exception.ConflictException;
import com.itsleonb.billsplittr.api.model.merchant.MerchantResponse;
import com.itsleonb.billsplittr.api.model.merchant.NewMerchantRequest;
import com.itsleonb.billsplittr.api.repository.merchant.MerchantRepository;
import com.itsleonb.billsplittr.api.service.MerchantService;
import com.itsleonb.billsplittr.impl.mapper.MerchantMapper;
import jakarta.validation.ConstraintViolation;
import jakarta.validation.ConstraintViolationException;
import jakarta.validation.Validator;
import lombok.AllArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;
import java.util.Set;

@Service
@AllArgsConstructor(onConstructor_ = @__(@Autowired))
public class MerchantServiceImpl implements MerchantService {
  private MerchantRepository merchantRepository;
  private Validator validator;

  @Override
  public MerchantResponse create(NewMerchantRequest request) {
    Set<ConstraintViolation<NewMerchantRequest>> constraintViolations = validator.validate(request);
    if (!constraintViolations.isEmpty()) {
      throw new ConstraintViolationException(constraintViolations);
    }

    Optional<Merchant> existingMerchant = merchantRepository.findByName(request.getName());
    if (existingMerchant.isPresent()) {
      throw new ConflictException(String.format("Merchant with name %s already exists", request.getName()));
    }

    Merchant merchantToCreate = MerchantMapper.fromNewRequest(request);
    Merchant createdMerchant = merchantRepository.save(merchantToCreate);

    return MerchantMapper.toResponse(createdMerchant);
  }

  @Override
  public List<MerchantResponse> find(String name) {
    List<Merchant> merchants = merchantRepository.searchByName(name);

    return MerchantMapper.toResponses(merchants);
  }
}
