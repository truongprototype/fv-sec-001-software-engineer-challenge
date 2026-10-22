> Đọc toàn bộ README.md của dự án.
>
> Chưa viết code và chưa tự tạo dữ liệu kết quả.
>
> Hãy thực hiện:
>
> 1. Trích xuất tất cả yêu cầu chức năng, kỹ thuật và bàn giao.
> 2. Liệt kê các trường hợp biên và điểm chưa được đặc tả rõ.
> 3. Đề xuất kiến trúc Go tối giản, dễ kiểm thử và xử lý CSV khoảng 1 GB bằng streaming.
> 4. Xác định thuật toán tổng hợp và tìm Top 10.
> 5. Đề xuất cách xử lý số tiền, chia cho 0, dữ liệu sai định dạng và các trường hợp đồng hạng.
> 6. Đưa ra kế hoạch test và benchmark có thể tái lập.
>
> Tạo IMPLEMENTATION_PLAN.md. Phân biệt yêu cầu bắt buộc với đề xuất của bạn. Không tự ý bổ sung tính năng ngoài phạm vi.

> Các đề xuất trong  IMPLEMENTATION_PLAN.md đã được duyệt, hãy tạo bộ khung dự án Go.
> Yêu cầu:
>
> - Dùng Go standard library nếu chưa có lý do rõ ràng để thêm dependency.
> - Tách CLI, logic aggregation và phần đọc/ghi CSV thành các package hoặc module hợp lý.
> - Có thể chạy bằng go run và build thành binary bằng go build.
> - Hỗ trợ --input và --output.
> - Xác thực tham số, xử lý lỗi và trả về exit code phù hợp.
> - Không triển khai toàn bộ thuật toán trong một file main.go.
> - Thêm .gitignore và hướng dẫn chạy cơ bản.
> Chỉ triển khai bộ khung, có thể chưa cần tối ưu hiệu năng. Chạy gofmt, go vet và go test ./...; báo rõ lệnh nào đã chạy và kết quả thực tế.

> Triển khai riêng phần "3A — CSV reader và aggregator" theo IMPLEMENTATION_PLAN.md.
>
> Điều kiện:
>
> - Chỉ sửa những file cần thiết cho phần này.
> - Viết unit test cho logic mới và các trường hợp biên.
> - Không thay đổi hành vi đã được kiểm thử ở các phần trước nếu không có lý do.
> - Không tối ưu vi mô khi chưa có benchmark.
> - Không thêm dependency nếu thư viện chuẩn đáp ứng yêu cầu.
> - Chạy gofmt, go vet và go test ./... sau khi hoàn thành.
> - Báo cáo file đã sửa, quyết định kỹ thuật, kết quả test và các vấn đề còn tồn tại.
>
> Chỉ báo cáo test thành công nếu thực sự chạy đúng.

> Triển khai riêng phần "3B — Metrics và ranking" theo IMPLEMENTATION_PLAN.md.
>
> Điều kiện:
>
> - Chỉ sửa những file cần thiết cho phần này.
> - Viết unit test cho logic mới và các trường hợp biên.
> - Không thay đổi hành vi đã được kiểm thử ở các phần trước nếu không có lý do.
> - Không tối ưu vi mô khi chưa có benchmark.
> - Không thêm dependency nếu thư viện chuẩn đáp ứng yêu cầu.
> - Chạy gofmt, go vet và go test ./... sau khi hoàn thành.
> - Báo cáo file đã sửa, quyết định kỹ thuật, kết quả test và các vấn đề còn tồn tại.
> Chỉ báo cáo test thành công nếu thực sự chạy đúng.

> Triển khai riêng phần "3C — CSV writer và CLI integration" theo IMPLEMENTATION_PLAN.md.
>
> Điều kiện:
>
> - Chỉ sửa những file cần thiết cho phần này.
> - Viết unit test cho logic mới và các trường hợp biên.
> - Không thay đổi hành vi đã được kiểm thử ở các phần trước nếu không có lý do.
> - Không tối ưu vi mô khi chưa có benchmark.
> - Không thêm dependency nếu thư viện chuẩn đáp ứng yêu cầu.
> - Chạy gofmt, go vet và go test ./... sau khi hoàn thành.
> - Báo cáo file đã sửa, quyết định kỹ thuật, kết quả test và các vấn đề còn tồn tại.
> Chỉ báo cáo test thành công nếu thực sự chạy đúng.

