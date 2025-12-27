package com.portfolio.blog.repository;

import java.util.Optional;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import com.portfolio.blog.model.BlogPost;

@Repository
public interface BlogRepository extends JpaRepository<BlogPost, Long> {
  Optional<BlogPost> findBySlug(String slug);

  boolean existsBySlug(String slug);

  Page<BlogPost> findByPublishedTrue(Pageable pageable);

  // Search title or tags
  Page<BlogPost> findByTitleContainingIgnoreCaseOrTagsContainingIgnoreCase(String title, String tags,
      Pageable pageable);
}
