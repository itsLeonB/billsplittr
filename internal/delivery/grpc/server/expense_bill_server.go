package server

import (
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr-protos/gen/go/expensebill/v1"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/ungerr"
	"github.com/rotisserie/eris"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExpenseBillServer struct {
	expensebill.UnimplementedExpenseBillServiceServer
	validate       *validator.Validate
	expenseBillSvc service.ExpenseBillService
}

func newExpenseBillServer(
	validate *validator.Validate,
	expenseBillSvc service.ExpenseBillService,
) expensebill.ExpenseBillServiceServer {
	return &ExpenseBillServer{
		validate:       validate,
		expenseBillSvc: expenseBillSvc,
	}
}

func (ebs *ExpenseBillServer) UploadStream(stream expensebill.ExpenseBillService_UploadStreamServer) error {
	var metadata *expensebill.BillMetadata
	var imageData []byte

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return eris.Wrap(err, "failed to receive stream data")
		}

		switch data := req.Data.(type) {
		case *expensebill.UploadStreamRequest_BillMetadata:
			if metadata != nil {
				return status.Error(codes.InvalidArgument, "metadata already received")
			}
			metadata = data.BillMetadata

		case *expensebill.UploadStreamRequest_Chunk:
			if metadata == nil {
				return ungerr.BadRequestError("metadata must be sent first")
			}

			nextSize := int64(len(imageData)) + int64(len(data.Chunk))
			if metadata.GetFileSize() > 0 && nextSize > metadata.GetFileSize() {
				return ungerr.BadRequestError("stream exceeds declared file size")
			}

			imageData = append(imageData, data.Chunk...)
		}
	}

	if metadata == nil {
		return ungerr.BadRequestError("no metadata received")
	}

	// Validate file size
	if int64(len(imageData)) != metadata.GetFileSize() {
		return ungerr.BadRequestError("actual file size doesn't match expected size")
	}

	// Parse UUIDs
	payerProfileID, err := ezutil.Parse[uuid.UUID](metadata.GetPayerProfileId())
	if err != nil {
		return err
	}

	creatorProfileID, err := ezutil.Parse[uuid.UUID](metadata.GetCreatorProfileId())
	if err != nil {
		return err
	}

	// Create domain request
	uploadReq := &dto.UploadBillRequest{
		PayerProfileID:   payerProfileID,
		CreatorProfileID: creatorProfileID,
		ImageData:        imageData,
		ContentType:      metadata.ContentType,
		Filename:         metadata.Filename,
		FileSize:         metadata.FileSize,
	}

	if err = ebs.validate.Struct(uploadReq); err != nil {
		return err
	}

	// Call service layer
	id, err := ebs.expenseBillSvc.Upload(stream.Context(), uploadReq)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&expensebill.UploadStreamResponse{Id: id.String()})
}
