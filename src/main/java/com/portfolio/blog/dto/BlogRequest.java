package com.portfolio.blog.dto;

import jakarta.validation.constraints.NotBlank;
import lombok.Data;

@Data
public class BlogRequest {
  @NotBlank(message = "Title is required")
  private String title;

  @NotBlank(message = "Content is required")
  private String content;

  private String coverImage;
  private String tags;
  private boolean published;
}
