package com.itsleonb.billsplittr.impl.service;

import com.itsleonb.billsplittr.api.entity.merchant.Merchant;
import com.itsleonb.billsplittr.api.entity.merchant.MerchantItem;
import com.itsleonb.billsplittr.api.exception.BadRequestException;
import com.itsleonb.billsplittr.api.exception.ConflictException;
import com.itsleonb.billsplittr.api.exception.NotFoundException;
import com.itsleonb.billsplittr.api.model.merchant.MerchantItemResponse;
import com.itsleonb.billsplittr.api.model.merchant.MerchantResponse;
import com.itsleonb.billsplittr.api.model.merchant.NewMerchantItemRequest;
import com.itsleonb.billsplittr.api.model.merchant.NewMerchantRequest;
import com.itsleonb.billsplittr.api.repository.merchant.MerchantItemRepository;
import com.itsleonb.billsplittr.api.repository.merchant.MerchantRepository;
import com.itsleonb.billsplittr.api.service.MerchantService;
import com.itsleonb.billsplittr.impl.mapper.MerchantMapper;
import jakarta.validation.ConstraintViolation;
import jakarta.validation.ConstraintViolationException;
import jakarta.validation.Validator;
import lombok.AllArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.UUID;

@Service
@AllArgsConstructor(onConstructor_ = @__(@Autowired))
public class MerchantServiceImpl implements MerchantService {
  private MerchantRepository merchantRepository;
  private MerchantItemRepository merchantItemRepository;
  private Validator validator;

  @Override
  public MerchantResponse create(NewMerchantRequest request) {
    validateRequest(request);

    if (merchantRepository.existsByName(request.getName())) {
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

  @Override
  public MerchantResponse getById(UUID id) {
    Optional<Merchant> merchant = merchantRepository.findWithItemsById(id);
    if (merchant.isEmpty()) {
      throw new NotFoundException(String.format("Merchant with ID: %s is not found", id));
    }

    return MerchantMapper.toResponse(merchant.get());
  }

  @Override
  public MerchantItemResponse createItem(UUID id, NewMerchantItemRequest request) {
    Optional<Merchant> merchant = merchantRepository.findById(id);
    if (merchant.isEmpty()) {
      throw new NotFoundException(String.format("Merchant with ID: %s is not found", id));
    }

    validateRequest(request);

    if (request.getPrice().compareTo(BigDecimal.ZERO) < 1) {
      throw new BadRequestException("Price must be greater than zero");
    }

    if (merchantItemRepository.existsByMerchantIdAndName(id, request.getName())) {
      throw new ConflictException(String.format(
        "Item with name %s already exists for merchant %s",
        request.getName(),
        id
      ));
    }

    MerchantItem merchantItem = MerchantMapper.fromNewItemRequest(merchant.get(), request);
    MerchantItem createdItem = merchantItemRepository.save(merchantItem);

    return MerchantMapper.toItemResponse(createdItem);
  }

  private void validateRequest(Object request) {
    Set<ConstraintViolation<Object>> constraintViolations = validator.validate(request);
    if (!constraintViolations.isEmpty()) {
      throw new ConstraintViolationException(constraintViolations);
    }
  }
}
