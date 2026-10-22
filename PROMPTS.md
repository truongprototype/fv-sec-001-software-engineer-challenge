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

> Hãy review repository như một reviewer độc lập. Đọc README.md và IMPLEMENTATION_PLAN.md, sau đó kiểm tra code và test hiện có.
>
> Ưu tiên tìm:
>
> 1. Sai logic aggregation, CTR, CPA hoặc ranking.
> 2. Lỗi khi impressions hoặc conversions bằng 0.
> 3. Sai lệch độ chính xác khi cộng spend.
> 4. Lỗi xử lý header, quoted CSV fields, dòng lỗi và file rỗng.
> 5. Bộ nhớ tăng theo số dòng thay vì số campaign.
> 6. Tài nguyên không được đóng, file output bị ghi dở hoặc lỗi bị nuốt.
> 7. Test thiếu hoặc test không thực sự kiểm chứng kết quả.
>
> Chưa sửa code. Trả về từng phát hiện gồm mức độ nghiêm trọng, vị trí, cách tái hiện và đề xuất sửa. Phân biệt lỗi đã xác minh với rủi ro giả định.

> - Spend và CPA được xuất với hai chữ số thập phân. Nếu input hợp lệ có độ chính xác vượt quá cent, phần hiển thị bị làm tròn; README cũng không nói không được tính trường hợp này, vậy thì hãy triển khai thêm.
> - Những giả định khác chưa cần làm.

