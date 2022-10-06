package internal

import (
	"context"
	api "github.com/jamalkaksouri/Book-Service-grpc-golang/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcServer struct {
	BookRepository BookRepository
	api.UnimplementedBookServiceServer
}

func NewGRPCServer(repository BookRepository) *grpc.Server {
	srv := grpcServer{
		BookRepository: repository,
	}

	gSrv := grpc.NewServer()
	api.RegisterBookServiceServer(gSrv, &srv)
	return gSrv
}

func (s *grpcServer) CreateBook(ctx context.Context, req *api.CreateBookRequest) (*api.CreateBookResponse, error) {
	book := &Book{
		Id:          0,
		Title:       req.Book.GetTitle(),
		Author:      req.Book.GetAuthor(),
		Description: req.Book.GetDescription(),
		Language:    req.Book.GetLanguage(),
		FinishTime:  req.Book.GetFinishTime().AsTime(),
	}
	id, err := s.BookRepository.CreateBook(ctx, book)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &api.CreateBookResponse{Id: int64(id)}, nil
}

func (s *grpcServer) RetrieveBook(ctx context.Context, req *api.RetrieveBookRequest) (*api.RetrieveBookResponse, error) {
	book, err := s.BookRepository.RetrieveBook(ctx, BookId(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	res := &api.RetrieveBookResponse{
		Book: &api.Book{
			Id:          int64(book.Id),
			Title:       book.Title,
			Author:      book.Author,
			Description: book.Description,
			Language:    book.Language,
			FinishTime:  timestamppb.New(book.FinishTime),
		},
	}
	return res, nil
}

func (s *grpcServer) UpdateBook(ctx context.Context, req *api.UpdateBookRequest) (*api.UpdateBookResponse, error) {
	book := &Book{
		Id:          BookId(req.Book.GetId()),
		Title:       req.Book.GetTitle(),
		Author:      req.Book.GetAuthor(),
		Description: req.Book.GetDescription(),
		Language:    req.Book.GetLanguage(),
		FinishTime:  req.Book.GetFinishTime().AsTime(),
	}
	err := s.BookRepository.UpdateBook(ctx, book)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &api.UpdateBookResponse{}, nil
}

func (s *grpcServer) DeleteBook(ctx context.Context, req *api.DeleteBookRequest) (*api.DeleteBookResponse, error) {
	err := s.BookRepository.DeleteBook(ctx, BookId(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &api.DeleteBookResponse{}, nil
}

func (s *grpcServer) ListBook(ctx context.Context, req *api.ListBookRequest) (*api.ListBookResponse, error) {
	books, _ := s.BookRepository.ListBook(ctx, req.Offset, req.Limit)
	res := &api.ListBookResponse{}
	var data []*api.Book
	for _, book := range books {
		b := &api.Book{
			Id:          int64(book.Id),
			Title:       book.Title,
			Author:      book.Author,
			Description: book.Description,
			Language:    book.Language,
			FinishTime:  timestamppb.New(book.FinishTime),
		}
		data = append(data, b)
	}
	res.Books = data
	return res, nil

}
